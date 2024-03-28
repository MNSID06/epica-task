import React, { useState } from "react";
import axios from "axios";

function App() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [getResponse, setGetResponse] = useState("");

  const handleSet = async () => {
    try {
      await axios.post("http://localhost:8080/set", { key, value });
    } catch (error) {
      console.error(error);
    }
  };

  const handleGet = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/get?key=${key}`);
      setGetResponse(response.data);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div>
      <input type="text" value={key} onChange={(e) => setKey(e.target.value)} />
      <input
        type="text"
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      <button onClick={handleSet}>Set</button>
      <button onClick={handleGet}>Get</button>
      <div>Get Response: {getResponse}</div>
    </div>
  );
}

export default App;
