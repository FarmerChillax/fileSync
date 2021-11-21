package core

// func StartServer(ctx context.Context, host, port string) (context.Context, error) {
// 	addr := net.JoinHostPort(host, port)
// 	listen, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 开启服务器
// 	fmt.Printf("Start service on: %v\n", addr)
// 	for {
// 		conn, err := listen.Accept()
// 		fmt.Printf("accept a connect from: %v\n", conn.RemoteAddr())
// 		if err != nil {
// 			log.Printf("accept failed, err: %v\n", err)
// 			continue
// 		}
// 		// 启动一个协程处理请求
// 		go GetLoaclSyncFile(conn, config.Path)
// 	}
// }

// func ServeMux(conn net.Conn) {
// 	defer conn.Close()
// 	for {
// 		buf := make([]byte, 4096)
// 		header := NewHeader()
// 		// 1. 解析传输过来的操作
// 		rn, err := conn.Read(buf)
// 		fmt.Println(buf[:rn])
// 		if err != nil {
// 			if err == io.EOF {
// 				log.Printf("Client %v offline.", conn.RemoteAddr())
// 				break
// 			}
// 			log.Printf("ServerMux read buffer err: %v\n", err)
// 			break
// 		}
// 		fmt.Printf("buf: %v\nbuf read size: %v\n", buf[:rn], rn)
// 		fmt.Printf("buf len: %v\nbuf cap: %v\n", len(buf), cap(buf))
// 		err = StructDecode(buf[:rn], &header)
// 		fmt.Println(header)
// 		if err != nil {
// 			log.Printf("Struct decode err: %v\n", err)
// 			break
// 		}
// 		// 2. 根据操作选择
// 		switch header.Type {
// 		case 1:
// 			fmt.Println("pull")
// 			fileEntry, _ := NewFileEntry("D:/projectCode/code.zip")
// 			fmt.Println(fileEntry)
// 			err := fileEntry.send(conn)
// 			if err != nil {
// 				fmt.Printf("发送的时候出错了: %v\n", err)
// 			}
// 		case 2:
// 			fmt.Println("push")
// 			conn.Write([]byte("server Push"))

// 		default:
// 			conn.Write([]byte("Method Not Allowed"))
// 		}
// 	}
// }
