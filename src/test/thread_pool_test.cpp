#include<boost/asio.hpp>
#include<iostream>
#include<mutex>
#include<chrono>
#include"../server/thread_pool.hpp"

using namespace std::literals;

int main(){
	cyg::asio_thread_pool pool;
	boost::asio::io_context::strand strand{pool.get_io_context()};
	boost::asio::steady_timer timer1{pool.get_io_context(),1s},
	timer2{pool.get_io_context(),500ms};
	timer1.async_wait(strand.wrap([](const boost::system::error_code&){
		std::cout<<"timer1 has finished"<<std::endl;
	}));
	timer2.async_wait(strand.wrap([](const boost::system::error_code&){
		std::cout<<"timer2 has finished"<<std::endl;
	}));
	pool.stop();

	return 0;
}
