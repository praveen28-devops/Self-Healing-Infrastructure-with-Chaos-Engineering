// Wait for the document to be fully loaded
document.addEventListener("DOMContentLoaded", () => {
    
    // Find all vote buttons
    const voteButtons = document.querySelectorAll(".vote-button");

    // Add a click event listener to each button
    voteButtons.forEach(button => {
        button.addEventListener("click", () => {
            const companyName = button.dataset.company;
            sendVote(companyName);
        });
    });
});

// Function to send the vote to the backend
async function sendVote(company) {
    try {
        const response = await fetch("/vote", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ company: company }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const result = await response.json();
        
        // Update the vote count on the page
        updateVoteCount(result.company, result.new_count);

    } catch (error) {
        console.error("Error sending vote:", error);
    }
}

// Function to update the count on the page
function updateVoteCount(company, newCount) {
    const countElement = document.getElementById(`count-${company}`);
    if (countElement) {
        // Add a simple "flash" animation
        countElement.style.transition = 'none';
        countElement.style.transform = 'scale(1.2)';
        countElement.style.opacity = '0.5';
        
        // Set the new number
        countElement.innerText = newCount;
        
        // Trigger the animation back to normal
        setTimeout(() => {
            countElement.style.transition = 'transform 0.2s ease, opacity 0.2s ease';
            countElement.style.transform = 'scale(1)';
            countElement.style.opacity = '1';
        }, 50); // A small delay to ensure the browser registers the change
    }
}