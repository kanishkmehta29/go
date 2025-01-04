import { useState } from "react";
import axios from "axios";
import "./App.css";

function App() {
  console.log("reload...\n")
  const [piclink, Setpiclink] = useState("#");
  const [querylink, Setquerylink] = useState("");

  async function getpic() {
    var res = await axios.get(querylink)
    Setpiclink(res.data.link)
  }

  return (
    <div>
      <form>
        <p>
          <input
            placeholder="Find your dreams ..."
            onChange={(e) => {
              Setquerylink("http://localhost:8080/search/" + e.target.value)
            }}
          ></input>
        </p>
        <p>
          <button onClick={getpic}>Search</button>
        </p>
      </form>
      <img src={piclink} alt="Search result" />
    </div>
  );
}

export default App;
