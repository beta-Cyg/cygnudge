#include<random>
#include<fstream>

int main(int argc,char **args){
	std::random_device r;
	std::default_random_engine e(r());
	std::uniform_int_distribution<int> uniform_dist(-0xffff,0xffff);
	std::ofstream in_s(args[1]),out_s(args[2]);
	int a{uniform_dist(e)},b{uniform_dist(e)};
	in_s<<a<<' '<<b;
	out_s<<a+b;

	return 0;
}
