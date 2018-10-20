cmds = {
  "ports": -> { Ports.print({"8000-8999", "3000-3999"}) },
  "cpu":   -> { CPU.print 3 },
}

if ENV["TMUXSTATUS_SPEC"]? != "1"
  if ARGV.size == 1 && (cmd = cmds[ARGV[0]]?)
    cmd
  else
    raise "usage: #{PROGRAM_NAME} #{cmds.keys.join "|"}"
  end.call
end

module Ports
  def self.print(ranges)
    puts open(ranges).sort.join(", ")
  end

  private def self.open(ranges)
    args = %w(-l -P -n -s TCP:LISTEN) + ranges.flat_map { |r| ["-i", "TCP:"+r] }
    proc = Process.new "lsof", args: args,
      output: Process::Redirect::Pipe,
      error: Process::Redirect::Pipe

    begin
      outp, err = proc.output.gets_to_end, proc.error.gets_to_end
    ensure
      st = proc.wait
    end
    st.success? || (st.exit_status == 256 && err == "") \
      || raise "lsof failed"

    EnumUtils.
      grep(outp.split("\n"), /.*:(\d+) \(LISTEN\)/) { |m,| m[1].to_i }.
      uniq
  end
end

lib LibC
  fun getuid : UInt
end

module CPU
  def self.print(n, len = 7)
    puts top.
      sort_by { |l| -l.usage }[0,n].
      map { |line| "%.1f %s" % {line.usage, truncate(short(line.cmd), len)} }.
      join(", ")
  end

  private def self.top
    lines = `ps -U #{LibC.getuid} -e -o pid,%cpu,comm`.split "\n"
    $?.success? || raise "ps failed"
    EnumUtils.grep(lines, /^\s*(\d+)\s+(\d+\.\d+)\s+(\S+)\s*$/) do |m,|
      TopLine.new pid: m[1].to_u32, usage: m[2].to_f32, cmd: m[3]
    end
  end

  protected def self.truncate(s, len)
    s.size <= len ? s : s[0,len]
  end

  protected def self.short(s) : String
    parts = s.split(" ", 2)
    if (path = parts[0]) && path.includes?('/')
      parts[0] = File.basename path
    end
    parts.join " "
  end

  class TopLine
    def initialize(@pid, @usage, @cmd)
    end

    getter pid : UInt32, usage : Float32, cmd : String
  end
end

module EnumUtils
  def self.grep(iter, re, &block : Regex::MatchData -> U) forall U
    ([] of U).tap do |res|
      iter.each do |val|
        if re === val
          res << yield($~, val)
        end
      end
    end
  end
end
