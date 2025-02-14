document.addEventListener("DOMContentLoaded", function () {
    const emailInput = document.getElementById("email");

    if (emailInput) {
        const email = localStorage.getItem("user_email");
        if (email) {
            emailInput.value = email;
        } else {
            alert("User email not found. Please register first.");
            window.location.href = "Sign Up.html";
        }
    } else {
        console.error("Error: Element with ID 'email' not found in DOM.");
    }

    // ✅ ตรวจสอบว่า form มีอยู่จริง
    const logisticsForm = document.getElementById("logistics-form");
    if (!logisticsForm) {
        console.error("❌ Error: Form with ID 'logistics-form' not found.");
        return;
    }

    logisticsForm.addEventListener("submit", async function (event) {
        console.log("✅ Form submit intercepted!");
        event.preventDefault(); // ✅ ป้องกันการโหลดหน้าใหม่

        const userID = localStorage.getItem("user_id");
        if (!userID) {
            alert("User ID not found. Please register first.");
            return;
        }

        const certificationCID = localStorage.getItem("certification_cid");
        console.log("✅ Certification CID from localStorage:", certificationCID);

        if (!certificationCID) {
            alert("No certification CID found. Please upload a certification first.");
            return;
        }

        // ✅ ฟังก์ชันช่วยให้ไม่เกิด null
        const getValue = (id) => {
            const element = document.getElementById(id);
            return element ? element.value : null;
        };

        const logisticsData = {
            userid: userID,
            company_name: getValue("company_name"),
            firstname: getValue("firstname"),
            lastname: getValue("lastname"),
            email: getValue("email"),
            address: getValue("address"),
            address2: getValue("address2"),
            areacode: getValue("areacode"),
            phone: getValue("phone"),
            post: getValue("post"),
            city: getValue("city"), 
            province: getValue("province"),
            country: getValue("country"),
            lineid: getValue("lineid"),
            facebook: getValue("facebook"),
            location_link: getValue("location_link")
        };

        try {
            console.log("📌 Sending logistics data to backend...");
            const response = await fetch("http://127.0.0.1:8080/api/v1/logistics", { 
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(logisticsData)
            });

            const textResponse = await response.text();
            console.log("📌 Server Response:", textResponse);

            if (!response.ok) {
                alert("Error: " + textResponse);
                return;
            }

            const data = JSON.parse(textResponse);
            console.log("📌 Data from Backend:", data); 

            // ✅ ตรวจสอบ logistics ID
            const logisticsID = data.logisticsid || data.logistics_id;
            if (!logisticsID) {
                alert("❌ Error: Logistics ID is missing from response!");
                return;
            }

            alert("✅ Logistics provider information saved successfully!");

            // ✅ บันทึกใบเซอร์หลังจากลงทะเบียน
            console.log("📌 Sending certification data to backend...");
            const certResponse = await fetch("http://127.0.0.1:8080/api/v1/certifications/create", { 
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    entity_type: "Logistics",  
                    entity_id: logisticsID,
                    certification_cid: certificationCID,
                    issued_date: "2025-02-07",
                    expiry_date: "2026-06-01"
                })
            });

            const certTextResponse = await certResponse.text();
            console.log("📌 Certification Server Response:", certTextResponse);

            if (!certResponse.ok) {
                console.error("❌ Failed to save certification:", certTextResponse);
                alert("Failed to save certification: " + certTextResponse);
                return;
            }

            alert("✅ Certification saved successfully!");
            window.location.href = "index.html";

        } catch (error) {
            console.error("❌ Error:", error);
            alert("An error occurred while saving data.");
        }
    });
});
