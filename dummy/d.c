
#include <stdio.h>

#define OB_VAR(x)  x
#define OB_FUNC
#define OB_CODE(...)


int OB_VAR(newvar)=2;

int OB_FUNC PrintFunc()
{
	int a,b,c;
	a = 0;
	b = 0;
	c = 0;
	OB_CODE(a,b,c);
	printf("hello world %d %d %d\n",a,b,c);
	return 0;
}

int main(int argc,char* argv[])
{
	newvar = 0;
	PrintFunc();
	return 0;
}