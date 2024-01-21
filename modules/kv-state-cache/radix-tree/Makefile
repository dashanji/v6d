DEBUG?= -g -ggdb
CFLAGS?= -O2 -Wall -W -std=c++11
LDFLAGS= -lm

# Uncomment the following two lines for coverage testing
#
# CFLAGS+=-fprofile-arcs -ftest-coverage
# LDFLAGS+=-lgcov

all: radix-test 

radix.o: radix.h
radix-test.o: radix.h

radix-test: radix-test.o radix.o
	$(CXX) -o $@ $^ $(LDFLAGS) $(DEBUG)

.c.o:
	$(CXX) -c $(CFLAGS) $(DEBUG) $<

clean:
	rm -f radix-test *.gcda *.gcov *.gcno *.o
