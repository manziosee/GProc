package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gproc/pkg/types"
)

type Dashboard struct {
	manager ProcessManager
}

type ProcessManager interface {
	List() []*types.Process
	Start(proc *types.Process) error
	Stop(id string) error
	Restart(id string) error
}

func NewDashboard(manager ProcessManager) *Dashboard {
	return &Dashboard{manager: manager}
}

func (d *Dashboard) Start(port int) error {
	http.HandleFunc("/", d.indexHandler)
	http.HandleFunc("/api/processes", d.processesHandler)
	http.HandleFunc("/api/stop", d.stopHandler)
	http.HandleFunc("/api/restart", d.restartHandler)

	fmt.Printf("Web dashboard available at http://localhost:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (d *Dashboard) indexHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>GProc Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .btn { padding: 5px 10px; margin: 2px; cursor: pointer; border: none; }
        .stop { background-color: #f44336; color: white; }
        .restart { background-color: #ff9800; color: white; }
    </style>
</head>
<body>
    <h1>GProc Process Manager</h1>
    <div id="processes"></div>
    <script>
        function loadProcesses() {
            fetch('/api/processes')
                .then(response => response.json())
                .then(data => {
                    const html = '<table><tr><th>Name</th><th>Status</th><th>PID</th><th>Restarts</th><th>Actions</th></tr>' +
                        data.map(p => '<tr><td>' + p.name + '</td><td>' + p.status + '</td><td>' + p.pid + '</td><td>' + p.restarts + '</td><td>' +
                        '<button class="btn stop" onclick="stopProcess(\'' + p.id + '\')">Stop</button>' +
                        '<button class="btn restart" onclick="restartProcess(\'' + p.id + '\')">Restart</button></td></tr>').join('') +
                        '</table>';
                    document.getElementById('processes').innerHTML = html;
                });
        }
        
        function stopProcess(id) {
            fetch('/api/stop?id=' + id, {method: 'POST'})
                .then(() => loadProcesses());
        }
        
        function restartProcess(id) {
            fetch('/api/restart?id=' + id, {method: 'POST'})
                .then(() => loadProcesses());
        }
        
        loadProcesses();
        setInterval(loadProcesses, 5000);
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (d *Dashboard) processesHandler(w http.ResponseWriter, r *http.Request) {
	processes := d.manager.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processes)
}

func (d *Dashboard) stopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	
	if err := d.manager.Stop(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (d *Dashboard) restartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	
	if err := d.manager.Restart(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}