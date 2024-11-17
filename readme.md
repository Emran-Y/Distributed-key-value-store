### Reflection


### **Concurrency Implementation**

The server leverages goroutines to handle multiple clients concurrently. Each client connection is processed in its own goroutine, ensuring requests from one client do not block others. The main server loop listens for incoming connections and spawns a new goroutine for each connection using `go handleClient(conn)`. This design enables efficient management of multiple clients simultaneously.

---

### **Challenges and Solutions**

1. **Race Conditions**  
   Concurrent access to a shared `map[string]string` initially caused race conditions when multiple clients attempted to read or modify the same key. To resolve this, a `sync.Mutex` was introduced to synchronize access to the shared map. Operations such as PUT, GET, DELETE, and LIST lock the mutex before accessing the map and unlock it afterward. Using `defer mu.Unlock()` ensures the mutex is released, even in case of errors.

2. **Client Disconnections**  
   Unexpected client disconnections led to errors during read and write operations. To handle this, error checks were added for every read and write. If a client disconnected, the connection was gracefully closed, and the associated goroutine terminated without impacting the server’s stability.

3. **Concurrency Overhead**  
   Although goroutines are lightweight, an unlimited number of connections could strain resources. While the current implementation supports the existing workload, introducing connection limits or a goroutine pool will be necessary for scalability in production environments.

4. **Deadlocks**  
   Incorrect handling of the mutex resulted in deadlocks, especially when errors caused early returns. To prevent this, `defer` was consistently used to ensure the mutex is unlocked regardless of how the function exits.

---

### **Maintaining Data Consistency**

Data consistency was ensured by synchronizing access to the shared store with `sync.Mutex`. This synchronization allowed operations to be atomic, avoiding issues like partial writes or stale reads. Although it introduced some contention under high load, it guaranteed correctness. Additionally, invalid operations like attempting `GET` or `DELETE` on non-existent keys returned appropriate error messages, preventing corruption or unexpected behavior.

---

### **Summary**

The server employs goroutines for concurrent client handling and uses a `sync.Mutex` to maintain data consistency and prevent race conditions. Challenges such as race conditions, client disconnections, and deadlocks were mitigated with thoughtful design and robust error handling. While the current implementation is stable and performant for the current scale, future scalability improvements like connection limits or goroutine pooling will be required for handling larger workloads.
