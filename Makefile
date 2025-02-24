buildall:
	make clean
	go build -o bin/ .
	cp temlate.xlsx bin/
	make buildclient

buildclient:
	cd client && npm install && npm run build
	mkdir -p bin/client/dist
	cp -r client/dist/* bin/client/dist

clean:
	rm -rf bin/
	rm -rf client/dist

publish:
	docker build -t lutzpfannenschmidt/anrechnungsstundenberechner . 
	docker push lutzpfannenschmidt/anrechnungsstundenberechner

docker_test:
	docker build -t test .
	make docker_run

docker_run:
	docker run --volume anrech_data:/app/pb_data -p 1337:80 test