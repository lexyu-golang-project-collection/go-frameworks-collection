import './style.css';

import { ProcessJSON } from '../wailsjs/go/main/App';

const inputElement = document.getElementById("jsonInput") as HTMLTextAreaElement;
const outputElement = document.getElementById("jsonOutput") as HTMLPreElement;

// Setup the greet function
window.greet = async function () {
    if (!inputElement || !outputElement) {
        console.error("Input or output element not found!");
        return;
    }

    const jsonString = inputElement.value;

    try {
        // 調用後端 ProcessJSON 方法
        const formattedJSON = await ProcessJSON(jsonString);
        console.log(formattedJSON)
        // 顯示格式化後的 JSON
        outputElement.textContent = formattedJSON;
    } catch (error) {
        // 處理錯誤
        console.error("Error processing JSON:", error);
        outputElement.textContent = "Invalid JSON or processing error: " + error;
    }
};

declare global {
    interface Window {
        greet: () => void;
    }
}
