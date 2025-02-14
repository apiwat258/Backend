document.getElementById("upload_certification").addEventListener("change", async function (event) {
    const file = event.target.files[0];
    if (!file) {
        alert("Please select a file.");
        return;
    }

    const formData = new FormData();
    formData.append("file", file);

    // ✅ ดึง Token (ถ้ามี)
    const token = localStorage.getItem("auth_token");

    try {
        console.log("📌 Uploading file to IPFS...");
        const response = await fetch("http://127.0.0.1:8080/api/v1/certifications/upload", { 
            method: "POST",
            headers: token ? { "Authorization": `Bearer ${token}` } : {}, 
            body: formData
        });

        const result = await response.json();
        console.log("✅ IPFS Upload Response:", result); // ✅ เพิ่ม Debug Log

        if (!response.ok || !result.cid) {
            alert("❌ Failed to upload file to IPFS");
            return;
        }

        const certificationCID = result.cid;  
        console.log("✅ Certification CID:", certificationCID);

        // ✅ เก็บ `CID` ไว้ใน localStorage
        localStorage.setItem("certification_cid", certificationCID);
        alert("File uploaded successfully! CID: " + certificationCID);

    } catch (error) {
        console.error("❌ Error uploading file:", error);
        alert("An error occurred while uploading.");
    }
});
