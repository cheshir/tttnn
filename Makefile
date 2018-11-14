init:
	docker build -t local/tf:latest .
	docker run -ti --name tttnn -p 9010:8888 -p 8536:8536 -p 0.0.0.0:9011:6006 -v `pwd`/.:/go/src/github.com/cheshir/tttnn local/tf
run:
	docker start tttnn
stop:
	docker stop -t 0 tttnn
delete:
	docker stop -t 0 tttnn
	docker rm tttnn
shell:
	docker exec -ti tttnn bash
