P=isrc
LINK_IP=127.0.0.1

$(P).so:
	go build --buildmode=plugin -o $@

.PHONY: upload
upload: $(P).so
	curl \
		-F "model.so=@./$(P).so" \
		$(LINK_IP):1372/api/models

.PHONY: clean
clean:
	rm $(P).so
