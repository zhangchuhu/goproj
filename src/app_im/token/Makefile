PROGRAM = libtoken_validate.a

AR = ar rcs

SRC_COMM := $(wildcard *.cpp) \
	$(wildcard ycloud_token/*.cpp) \
	$(wildcard ycloud_token/encrypt/*.cpp) 
OBJ_COMM := $(SRC_COMM:.cpp=.o)

CXXFLAG = -g 
INCLUDE := -I./ycloud_token/include \
	-I ./ycloud_token/encrypt/include \

.SUFFIXES: .o .cpp .c
.cpp.o:
	$(CXX) $(CXXFLAG) $(INCLUDE) -c -o $@ $<

all: $(PROGRAM)

$(PROGRAM): $(OBJ_COMM)
	$(AR)  $@  $(OBJ_COMM) 

clean: 
	rm -f $(PROGRAM) $(OBJ_COMM)
