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

    // ✅ เมื่อกด Submit จะส่งข้อมูลโรงงาน + CID ที่อัปโหลดไว้ไปยัง Backend
    const factoryForm = document.getElementById("factory-form");
    if (factoryForm) {
        factoryForm.addEventListener("submit", async function (event) {
            event.preventDefault();

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

            const factoryData = {
                userid: userID,
                company_name: document.getElementById("company_name").value,
                firstname: document.getElementById("firstname").value,
                lastname: document.getElementById("lastname").value,
                email: document.getElementById("email").value,
                address: document.getElementById("address").value,
                address2: document.getElementById("address2").value,
                areacode: document.getElementById("areacode").value,
                phone: document.getElementById("phone").value,
                post: document.getElementById("post").value,
                city: document.getElementById("city").value,
                location_link: document.getElementById("location_link").value
            };

            try {
                const response = await fetch("http://127.0.0.1:8080/api/v1/factories", { 
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(factoryData)
                });

                const textResponse = await response.text();
                console.log("📌 Server Response:", textResponse); 

                if (!response.ok) {
                    alert("Error: " + textResponse);
                    return;
                }

                const data = JSON.parse(textResponse);
                alert("Factory information saved successfully!");

                // ✅ บันทึกใบเซอร์หลังจากลงทะเบียนโรงงาน
                console.log("📌 Sending certification data to backend...");
                const certResponse = await fetch("http://127.0.0.1:8080/api/v1/certifications/create", { 
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        entity_type: "Factory",
                        entity_id: data.factory_id,
                        certification_cid: certificationCID,
                        issued_date: "2025-02-07",
                        expiry_date: "2026-06-01"
                    })
                });

                const certData = await certResponse.json();
                console.log("✅ Certification Response:", certData);

                if (certResponse.ok) {
                    alert("Certification saved successfully!");
                } else {
                    console.error("❌ Failed to save certification:", certData.error);
                    alert("Failed to save certification: " + certData.error);
                }

                window.location.href = "index.html";

            } catch (error) {
                console.error("❌ Error:", error);
                alert("An error occurred while saving data.");
            }
        });
    } else {
        console.error("Error: Form with ID 'factory-form' not found in DOM.");
    }
});
