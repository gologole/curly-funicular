run_elk:
	cd ELK
	sudo docker-compose up 

run_auth:
	go run jwt/main.go

run_authtest:
	grpcurl -plaintext -import-path ./ -proto auth.proto -d '{"username": "user", "password": "testpassword"}' localhost:50051 auth.AuthService/Login


run_elktest:
	chmod +x tests.sh
	./tests.sh


