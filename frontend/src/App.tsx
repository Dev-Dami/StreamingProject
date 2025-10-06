import React, { useEffect, useRef, useState } from "react";

export default function StreamPage() {
  const [status, setStatus] = useState("Connecting...");
  const [isConnected, setIsConnected] = useState(false);
  const [frameCount, setFrameCount] = useState(0);
  const imgRef = useRef<HTMLImageElement>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const connectWebSocket = () => {
      const ws = new WebSocket("ws://localhost:8080/ws");
      ws.binaryType = "arraybuffer";
      wsRef.current = ws;

      ws.onopen = () => {
        setStatus("Connected");
        setIsConnected(true);
        console.log("WebSocket connected");
      };

      ws.onclose = () => {
        setStatus("Disconnected");
        setIsConnected(false);
        console.log("WebSocket disconnected");
        // Attempt to reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
      };

      ws.onerror = (error) => {
        setStatus("Connection Error");
        setIsConnected(false);
        console.error("WebSocket error:", error);
      };

      ws.onmessage = (event) => {
        if (typeof event.data !== "object") return;
        
        const blob = new Blob([event.data], { type: "image/jpeg" });
        const url = URL.createObjectURL(blob);

        if (imgRef.current) {
          // Clean up previous object URL to prevent memory leaks
          if (imgRef.current.src.startsWith("blob:")) {
            URL.revokeObjectURL(imgRef.current.src);
          }
          imgRef.current.src = url;
          setFrameCount(prev => prev + 1);
        }
      };
    };

    connectWebSocket();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  const handleReconnect = () => {
    if (wsRef.current) {
      wsRef.current.close();
    }
    setStatus("Reconnecting...");
    setFrameCount(0);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
      {/* Header */}
      <header className="bg-gray-800/50 backdrop-blur-sm border-b border-gray-700/50 px-6 py-4">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">VS</span>
            </div>
            <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
              Video Streamer
            </h1>
          </div>
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2">
              <div className={`w-3 h-3 rounded-full ${
                isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
              }`} />
              <span className={`text-sm font-medium ${
                isConnected ? 'text-green-400' : 'text-red-400'
              }`}>
                {status}
              </span>
            </div>
            {!isConnected && (
              <button
                onClick={handleReconnect}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-sm font-medium transition-colors duration-200"
              >
                Reconnect
              </button>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Video Stream */}
          <div className="lg:col-span-3">
            <div className="bg-gray-800/30 backdrop-blur-sm rounded-2xl border border-gray-700/50 overflow-hidden shadow-2xl">
              <div className="bg-gray-800/50 px-6 py-4 border-b border-gray-700/50">
                <h2 className="text-lg font-semibold text-gray-200">Live Video Feed</h2>
              </div>
              <div className="p-6">
                <div className="aspect-video bg-black rounded-xl overflow-hidden border border-gray-700/30 relative">
                  {!isConnected && (
                    <div className="absolute inset-0 flex items-center justify-center bg-gray-900/80">
                      <div className="text-center">
                        <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
                        <p className="text-gray-400 text-lg">
                          {status === "Connecting..." ? "Connecting to stream..." : "Connection lost"}
                        </p>
                        <p className="text-gray-500 text-sm mt-2">
                          Make sure the backend server is running
                        </p>
                      </div>
                    </div>
                  )}
                  <img
                    ref={imgRef}
                    alt="Video Stream"
                    className="w-full h-full object-contain"
                    style={{ display: isConnected ? 'block' : 'none' }}
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Statistics Panel */}
          <div className="lg:col-span-1 space-y-6">
            <div className="bg-gray-800/30 backdrop-blur-sm rounded-2xl border border-gray-700/50 p-6 shadow-xl">
              <h3 className="text-lg font-semibold text-gray-200 mb-4">Stream Statistics</h3>
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Status:</span>
                  <span className={`font-medium ${
                    isConnected ? 'text-green-400' : 'text-red-400'
                  }`}>
                    {status}
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Frames Received:</span>
                  <span className="font-medium text-blue-400">{frameCount.toLocaleString()}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Protocol:</span>
                  <span className="font-medium text-purple-400">WebSocket</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Format:</span>
                  <span className="font-medium text-orange-400">JPEG</span>
                </div>
              </div>
            </div>

            <div className="bg-gray-800/30 backdrop-blur-sm rounded-2xl border border-gray-700/50 p-6 shadow-xl">
              <h3 className="text-lg font-semibold text-gray-200 mb-4">Server Info</h3>
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Backend:</span>
                  <span className="font-medium text-gray-300">Go WebSocket</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Port:</span>
                  <span className="font-medium text-gray-300">8080</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Video Processing:</span>
                  <span className="font-medium text-gray-300">FFmpeg</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Frame Rate:</span>
                  <span className="font-medium text-gray-300">10 FPS</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
