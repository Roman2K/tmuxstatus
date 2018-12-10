cmds = {
  "ports": -> { Ports.print({8000..8999, 3000..3999}) },
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
    lines = `netstat -tln`.split "\n"
    $?.success? || raise "netstat failed"
    EnumUtils.
      grep(lines, /[\d:]:(\d+) /) { |m,| m[1].to_i }.
      select { |n| ranges.any? { |r| r.includes? n } }.
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
      map { |line| "%g %s" % {line.usage, truncate(short(line.cmd), len)} }.
      join(", ")
  end

  record TopLine, pid : UInt32, usage : Float32, cmd : String

  private def self.top
    lines = `ps -U #{LibC.getuid} -e -o pid,%cpu,comm`.split "\n"
    $?.success? || raise "ps failed"
    EnumUtils.grep(lines, /^\s*(\d+)\s+(\d+(?:\.\d+)?)\s+(\S+)\s*$/) do |m,|
      TopLine.new pid: m[1].to_u32, usage: m[2].to_f32, cmd: m[3]
    end
  end

  protected def self.truncate(s, len)
    s.size <= len ? s : s[0,len]
  end

  protected def self.short(s) : String
    parts = s.split(" ", 2)
    if (path = parts[0]).includes?('/')
      parts[0] = File.basename path
    end
    parts.join " "
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
