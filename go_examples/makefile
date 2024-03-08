include ../Makefile.common

.PHONY: all

ifdef LOCAL_INCLUDE
LOC_INC="-I$(LOCAL_INCLUDE)"
endif
ifdef LOCAL_LIB_ROOT
LOC_LIB="-L$(LOCAL_LIB_ROOT)"
endif

CXXFLAGS = $(CXX11FLAGS) $(OS_CXXFLAGS) -pthread -Wall -g -O2 -DSSL_LIB_INIT -I${ZROOT}/deps/root/include -I${ZROOT}/root/include $(LOC_INC)
LDFLAGS = -L${ZROOT}/deps/root/lib -L${ZROOT}/root/lib $(LOC_LIB)
LIBS = -lcrypto -lrelic -lrelic_ec -lopenabe

all: liblibrary.so

liblibrary.so: lib_bridge.cpp
	$(CXX) $(CXXFLAGS) $(LDFLAGS) -fPIC -shared -o liblibrary.so lib_bridge.cpp $(LIBS)