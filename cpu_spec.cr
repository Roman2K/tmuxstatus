require "spec"

ENV["TMUXSTATUS_SPEC"] = "1"
require "./main"

module CPUMethods
  include CPU

  macro expose(m)
    # TODO &block
    def self.{{m}}(*args, **opts)
      CPU.{{m}}(*args, **opts)
    end
  end

  expose truncate
  expose short
end

describe "CPU" do
  m = CPUMethods

  describe ".truncate" do
    describe "actual <= len" do
      it "doesn't truncate" do
        m.truncate("", 1).should eq ""
        m.truncate("x", 1).should eq "x"
      end
    end

    describe "actual > len" do
      it "truncates" do
        m.truncate("abc", 2).should eq "ab"
        m.truncate("abcd", 2).should eq "ab"
      end
    end
  end

  describe ".short" do
    describe "already short" do
      it "doesn't shorten" do
        m.short("myexe").should eq "myexe"
        m.short("myexe a").should eq "myexe a"
      end
    end

    describe "empty" do
      it "returns cmd as-is" do
        m.short("").should eq ""
      end
    end

    describe "exe contains path" do
      it "shortens" do
        m.short("mydir/myexe").should eq "myexe"
        m.short("mydir/myexe a").should eq "myexe a"
        m.short("/mydir/myexe").should eq "myexe"
        m.short("/mydir/myexe a").should eq "myexe a"
      end
    end
  end
end
