DEBUG?= -g -ggdb
CXXFLAGS?= -O2 -Wall  -W -std=c++11
#-DRAX_DEBUG_MSG
LDFLAGS= # -lm (Not needed for C++, but keep if other libraries need linking)
# Uncomment the following two lines for coverage testing
#
# CXXFLAGS+=-fprofile-arcs -ftest-coverage
# LDFLAGS+=-lgcov

all: radix-test 

radix.o: radix.h
radix-test.o: radix.h

#export DYLD_LIBRARY_PATH="$DYLD_LIBRARY_PATH:/usr/local/lib" && g++ -o $@ $^ ../lz4/lib/*.o $(LDFLAGS) $(DEBUG)
radix-test: radix-test.o radix.o
	g++ -o $@ $^ $(LDFLAGS) $(DEBUG)

.cpp.o:
	g++ -c $(CXXFLAGS) $(DEBUG) $<

clean:
	rm -f radix-test *.gcda *.gcov *.gcno *.o
