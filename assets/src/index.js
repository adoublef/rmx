const app = document.getElementById("app");

app.innerHTML = `
    <h1>Welcome, Home</h1>
    <p>To create a new session tap the button below</p>
    <button>Click Me!</button>

    <div id="session"></div>
`;

document.querySelector("button").addEventListener("click", async e => {
    try {
        const r = await fetch("/api/v1/jam/create");
        const { sessionId } = await r.json();

        session.id = sessionId;
    } catch (e) {
        console.error(e.message);
    }
});

const session = new Proxy({ id: "" }, {
    set(obj, prop, value) {
        let v = 0;
        switch (prop) {
            case "id":
                obj[prop] = value;
                document.getElementById("session").textContent = window.location.href + "play/" + value;
                return true;

            default: return false;
        }
    }
});