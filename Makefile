init:
	docker build -t local/tf:latest .
	docker run -ti --name 3tnn -p 9010:8888 -p 0.0.0.0:9011:6006 -v `pwd`/.:/app local/tf
run:
	docker start 3tnn
stop:
	docker stop -t 0 3tnn
delete:
	docker stop -t 0 3tnn
	docker rm 3tnn
shell:
	docker exec -ti 3tnn bash
