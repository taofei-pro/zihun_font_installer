import React, { useState } from "react";
import { SelectDirectory, InstallFonts } from "../wailsjs/go/main/App";
import "./App.css";

function App() {
  const [outputPath, setOutputPath] = useState("");

  // 选择目标路径
  const selectOutputPath = async () => {
    try {
      const result = await SelectDirectory();
      if (result) {
        setOutputPath(result);
      }
    } catch (err) {
      alert(`选择目录失败: ${err}`);
    }
  };

  // 安装字体
  const installFonts = async () => {
    if (!outputPath) {
      alert("请选择目标路径");
      return;
    }

    try {
      await InstallFonts(outputPath);
      alert("字体安装成功！");
    } catch (err) {
      alert(`安装失败: ${err}`);
    }
  };

  return (
    <div className="container">
      <h1 className="title">字体安装器</h1>
      <button onClick={selectOutputPath}>选择目标路径</button>
      <button onClick={installFonts}>安装字体</button>
      {outputPath && <p className="path">当前选择的路径: {outputPath}</p>}
    </div>
  );
}

export default App;
