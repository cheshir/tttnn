init:
	docker build -t local/tf:latest .
	docker run -ti --name tttnn -p 9010:8888 -p 0.0.0.0:9011:6006 -v `pwd`/.:/app local/tf
run:
	docker start tttnn
stop:
	docker stop -t 0 tttnn
delete:
	docker stop -t 0 tttnn
	docker rm tttnn
shell:
	docker exec -ti tttnn bash
