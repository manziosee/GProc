package sdks

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"gproc/pkg/types"
)

type SDKGenerator struct {
	config    *types.SDKConfig
	templates map[string]*template.Template
}

type SDKLanguage struct {
	Name       string `json:"name"`
	Extension  string `json:"extension"`
	Package    string `json:"package"`
	Repository string `json:"repository"`
}

func NewSDKGenerator(config *types.SDKConfig) *SDKGenerator {
	generator := &SDKGenerator{
		config:    config,
		templates: make(map[string]*template.Template),
	}
	
	generator.loadTemplates()
	return generator
}

func (s *SDKGenerator) GenerateAllSDKs(outputDir string) error {
	languages := []SDKLanguage{
		{Name: "nodejs", Extension: "js", Package: "gproc-client", Repository: "npm"},
		{Name: "python", Extension: "py", Package: "gproc-client", Repository: "pypi"},
		{Name: "java", Extension: "java", Package: "gproc-client", Repository: "maven"},
		{Name: "rust", Extension: "rs", Package: "gproc-client", Repository: "crates.io"},
	}
	
	for _, lang := range languages {
		if err := s.GenerateSDK(lang, outputDir); err != nil {
			return fmt.Errorf("failed to generate %s SDK: %v", lang.Name, err)
		}
	}
	
	return nil
}

func (s *SDKGenerator) GenerateSDK(lang SDKLanguage, outputDir string) error {
	langDir := filepath.Join(outputDir, "sdks", lang.Name)
	if err := os.MkdirAll(langDir, 0755); err != nil {
		return err
	}
	
	switch lang.Name {
	case "nodejs":
		return s.generateNodeJSSDK(langDir)
	case "python":
		return s.generatePythonSDK(langDir)
	case "java":
		return s.generateJavaSDK(langDir)
	case "rust":
		return s.generateRustSDK(langDir)
	default:
		return fmt.Errorf("unsupported language: %s", lang.Name)
	}
}

func (s *SDKGenerator) generateNodeJSSDK(outputDir string) error {
	// Generate package.json
	packageJSON := `{
  "name": "gproc-client",
  "version": "1.0.0",
  "description": "GProc Node.js SDK for process management",
  "main": "index.js",
  "scripts": {
    "test": "jest"
  },
  "dependencies": {
    "axios": "^1.0.0",
    "ws": "^8.0.0"
  },
  "keywords": ["gproc", "process", "management"],
  "author": "GProc Team",
  "license": "MIT"
}`
	
	if err := s.writeFile(filepath.Join(outputDir, "package.json"), packageJSON); err != nil {
		return err
	}
	
	// Generate main client
	clientJS := `const axios = require('axios');
const WebSocket = require('ws');

class GProcClient {
  constructor(endpoint, apiKey) {
    this.endpoint = endpoint;
    this.apiKey = apiKey;
    this.client = axios.create({
      baseURL: endpoint,
      headers: { 'Authorization': 'Bearer ' + apiKey }
    });
  }

  async startProcess(name, command, options = {}) {
    const response = await this.client.post('/api/v1/processes', {
      name, command, ...options
    });
    return response.data;
  }

  async stopProcess(name) {
    const response = await this.client.delete('/api/v1/processes/' + name);
    return response.data;
  }

  async listProcesses() {
    const response = await this.client.get('/api/v1/processes');
    return response.data;
  }

  async getProcessLogs(name, lines = 100) {
    const response = await this.client.get('/api/v1/processes/' + name + '/logs?lines=' + lines);
    return response.data;
  }

  connectWebSocket() {
    const ws = new WebSocket(this.endpoint.replace('http', 'ws') + '/ws');
    return ws;
  }
}

module.exports = GProcClient;`
	
	return s.writeFile(filepath.Join(outputDir, "index.js"), clientJS)
}

func (s *SDKGenerator) generatePythonSDK(outputDir string) error {
	// Generate setup.py
	setupPy := `from setuptools import setup, find_packages

setup(
    name="gproc-client",
    version="1.0.0",
    description="GProc Python SDK for process management",
    packages=find_packages(),
    install_requires=[
        "requests>=2.25.0",
        "websocket-client>=1.0.0"
    ],
    author="GProc Team",
    license="MIT",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)`
	
	if err := s.writeFile(filepath.Join(outputDir, "setup.py"), setupPy); err != nil {
		return err
	}
	
	// Create package directory
	pkgDir := filepath.Join(outputDir, "gproc_client")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return err
	}
	
	// Generate main client
	clientPy := `import requests
import websocket
import json

class GProcClient:
    def __init__(self, endpoint, api_key):
        self.endpoint = endpoint
        self.api_key = api_key
        self.session = requests.Session()
        self.session.headers.update({'Authorization': f'Bearer {api_key}'})

    def start_process(self, name, command, **options):
        data = {'name': name, 'command': command, **options}
        response = self.session.post(f'{self.endpoint}/api/v1/processes', json=data)
        response.raise_for_status()
        return response.json()

    def stop_process(self, name):
        response = self.session.delete(f'{self.endpoint}/api/v1/processes/{name}')
        response.raise_for_status()
        return response.json()

    def list_processes(self):
        response = self.session.get(f'{self.endpoint}/api/v1/processes')
        response.raise_for_status()
        return response.json()

    def get_process_logs(self, name, lines=100):
        response = self.session.get(f'{self.endpoint}/api/v1/processes/{name}/logs?lines={lines}')
        response.raise_for_status()
        return response.json()

    def connect_websocket(self):
        ws_url = self.endpoint.replace('http', 'ws') + '/ws'
        return websocket.WebSocket(ws_url)`
	
	if err := s.writeFile(filepath.Join(pkgDir, "__init__.py"), clientPy); err != nil {
		return err
	}
	
	return s.writeFile(filepath.Join(pkgDir, "__init__.py"), "from .client import GProcClient")
}

func (s *SDKGenerator) generateJavaSDK(outputDir string) error {
	// Generate pom.xml
	pomXML := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    
    <groupId>io.gproc</groupId>
    <artifactId>gproc-client</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>
    
    <name>GProc Java SDK</name>
    <description>GProc Java SDK for process management</description>
    
    <properties>
        <maven.compiler.source>11</maven.compiler.source>
        <maven.compiler.target>11</maven.compiler.target>
    </properties>
    
    <dependencies>
        <dependency>
            <groupId>com.squareup.okhttp3</groupId>
            <artifactId>okhttp</artifactId>
            <version>4.9.0</version>
        </dependency>
        <dependency>
            <groupId>com.fasterxml.jackson.core</groupId>
            <artifactId>jackson-databind</artifactId>
            <version>2.13.0</version>
        </dependency>
    </dependencies>
</project>`
	
	if err := s.writeFile(filepath.Join(outputDir, "pom.xml"), pomXML); err != nil {
		return err
	}
	
	// Create source directory
	srcDir := filepath.Join(outputDir, "src", "main", "java", "io", "gproc", "client")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return err
	}
	
	// Generate main client
	clientJava := `package io.gproc.client;

import okhttp3.*;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.util.Map;

public class GProcClient {
    private final String endpoint;
    private final String apiKey;
    private final OkHttpClient client;
    private final ObjectMapper mapper;

    public GProcClient(String endpoint, String apiKey) {
        this.endpoint = endpoint;
        this.apiKey = apiKey;
        this.client = new OkHttpClient();
        this.mapper = new ObjectMapper();
    }

    public Map<String, Object> startProcess(String name, String command, Map<String, Object> options) throws IOException {
        options.put("name", name);
        options.put("command", command);
        
        RequestBody body = RequestBody.create(
            mapper.writeValueAsString(options),
            MediaType.get("application/json")
        );
        
        Request request = new Request.Builder()
            .url(endpoint + "/api/v1/processes")
            .header("Authorization", "Bearer " + apiKey)
            .post(body)
            .build();
            
        try (Response response = client.newCall(request).execute()) {
            return mapper.readValue(response.body().string(), Map.class);
        }
    }

    public Map<String, Object> stopProcess(String name) throws IOException {
        Request request = new Request.Builder()
            .url(endpoint + "/api/v1/processes/" + name)
            .header("Authorization", "Bearer " + apiKey)
            .delete()
            .build();
            
        try (Response response = client.newCall(request).execute()) {
            return mapper.readValue(response.body().string(), Map.class);
        }
    }
}`
	
	return s.writeFile(filepath.Join(srcDir, "GProcClient.java"), clientJava)
}

func (s *SDKGenerator) generateRustSDK(outputDir string) error {
	// Generate Cargo.toml
	cargoToml := `[package]
name = "gproc-client"
version = "1.0.0"
edition = "2021"
description = "GProc Rust SDK for process management"
license = "MIT"
repository = "https://github.com/gproc/rust-sdk"

[dependencies]
reqwest = { version = "0.11", features = ["json"] }
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
tokio = { version = "1.0", features = ["full"] }
tokio-tungstenite = "0.17"`
	
	if err := s.writeFile(filepath.Join(outputDir, "Cargo.toml"), cargoToml); err != nil {
		return err
	}
	
	// Create src directory
	srcDir := filepath.Join(outputDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return err
	}
	
	// Generate main client
	clientRust := `use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct GProcClient {
    endpoint: String,
    api_key: String,
    client: Client,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProcessRequest {
    pub name: String,
    pub command: String,
    #[serde(flatten)]
    pub options: HashMap<String, serde_json::Value>,
}

impl GProcClient {
    pub fn new(endpoint: String, api_key: String) -> Self {
        Self {
            endpoint,
            api_key,
            client: Client::new(),
        }
    }

    pub async fn start_process(
        &self,
        name: &str,
        command: &str,
        options: HashMap<String, serde_json::Value>,
    ) -> Result<serde_json::Value, reqwest::Error> {
        let request = ProcessRequest {
            name: name.to_string(),
            command: command.to_string(),
            options,
        };

        let response = self
            .client
            .post(&format!("{}/api/v1/processes", self.endpoint))
            .header("Authorization", format!("Bearer {}", self.api_key))
            .json(&request)
            .send()
            .await?;

        response.json().await
    }

    pub async fn stop_process(&self, name: &str) -> Result<serde_json::Value, reqwest::Error> {
        let response = self
            .client
            .delete(&format!("{}/api/v1/processes/{}", self.endpoint, name))
            .header("Authorization", format!("Bearer {}", self.api_key))
            .send()
            .await?;

        response.json().await
    }

    pub async fn list_processes(&self) -> Result<serde_json::Value, reqwest::Error> {
        let response = self
            .client
            .get(&format!("{}/api/v1/processes", self.endpoint))
            .header("Authorization", format!("Bearer {}", self.api_key))
            .send()
            .await?;

        response.json().await
    }
}`
	
	return s.writeFile(filepath.Join(srcDir, "lib.rs"), clientRust)
}

func (s *SDKGenerator) loadTemplates() {
	// Load SDK templates (simplified for demo)
}

func (s *SDKGenerator) writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}