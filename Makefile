BINARIES = boomkat boomkat.arm
BLDDIR = build
SRCDIR = .

all : $(BINARIES)

$(BINARIES) : %: $(BLDDIR)/%

$(BLDDIR)/boomkat:
	mkdir -p $(BLDDIR)
	cd $(SRCDIR) && go build -o $(abspath $@)

$(BLDDIR)/boomkat.arm:
	mkdir -p $(BLDDIR)
	cd $(SRCDIR) && GOOS=linux GOARCH=arm go build -o $(abspath $@)

clean:
	rm -fr $(BLDDIR)
