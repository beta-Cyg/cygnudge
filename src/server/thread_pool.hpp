#ifndef CYGNUDGE_THREAD_POOL_HPP
#define CYGNUDGE_THREAD_POOL_HPP

#include<boost/asio.hpp>

#include<vector>
#include<memory>

namespace cyg{
	class asio_thread_pool{
	public:
		asio_thread_pool(std::size_t size=std::thread::hardware_concurrency()):
		work_(new boost::asio::io_context::work(io_context_)){
			for(std::size_t i=0;i<size;++i){
				threads_.emplace_back([this](){
					io_context_.run();
				});
			}
		}
	
		asio_thread_pool(const asio_thread_pool&)=delete;
		asio_thread_pool& operator=(const asio_thread_pool&)=delete;

		boost::asio::io_context& get_io_context(){
			return io_context_;
		}
	
		void stop(){
			work_.reset();
			for(auto &t:threads_){
				t.join();
			}
		}
	private:
		boost::asio::io_context io_context_;
		std::unique_ptr<boost::asio::io_context::work> work_;
		std::vector<std::thread> threads_;
	};
}

#endif
