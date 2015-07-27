
void OB_FUNC print_hello()
{
	printf("hello world\n");
	return;
}

int OB_VAR varcode;

int main(int argc,char* argv[])
{
	OB_CODE(argc,varcode);
	print_hello();
}

