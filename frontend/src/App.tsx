import React, { useEffect, useRef, useState } from "react";

export default function StreamPage() {
  const [status, setStatus] = useState("Connecting...");
  const imgRef = useRef(null);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.binaryType = "arraybuffer";

    ws.onopen = () => setStatus("Connected");
    ws.onclose = () => setStatus("Disconnected");
    ws.onerror = () => setStatus("Error");

    ws.onmessage = (event) => {
      if (typeof event.data !== "object") return;
      const blob = new Blob([event.data], { type: "image/jpeg" });
      const url = URL.createObjectURL(blob);

      // update <img> source
      if (imgRef.current) {
        imgRef.current.src = url;
      }
    };

    return () => ws.close();
  }, []);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-950 text-white">
      <h1 className="text-2xl font-bold mb-4">Live Stream</h1>
      <p className="mb-2 text-gray-400">{status}</p>
      <div className="rounded-xl overflow-hidden border border-gray-700 shadow-lg">
        <img
          ref={imgRef}
          alt="Video Stream"
          className="max-w-full max-h-[80vh] object-contain bg-black"
        />
      </div>
    </div>
  );
}
