import EntryBar from "./components/EntryBar";
import "./App.css";
import { useEffect, useState } from "react";
import axios from "axios";
import { EntryProps } from "./components/EntryBar";
import CustomButton from "./components/CustomButton";
import AddEntryForm from "./components/AddEntryForm";
import ChangeEntryForm from "./components/ChangeEntryForm";

function App() {
  const [entries, setEntries] = useState([]);
  const [addNewEntry,setAddNewEntry] = useState(false);
  const [changeEntry,setChangeEntry] = useState(false);
  const [reload,setReload] = useState(false);
  const [currentId,setCurrentId] = useState('')

  useEffect(() => {
    axios.get("http://localhost:8080/entry").then((res) => {
      setEntries(res.data.entries);
      console.log(res.data.entries)
    })
    .catch(err => console.error("Error fetching entries:", err))
  }, [addNewEntry,reload,changeEntry]);
  
  if(addNewEntry){
    return(
      <div className="flex justify-center items-center h-screen">
        <div>
          <AddEntryForm setAddNewEntry={setAddNewEntry}/>
        </div>
      </div>
    )
  }

  if(changeEntry){
    return(
      <div className="flex justify-center items-center h-screen">
        <div>
          <ChangeEntryForm setChangeEntry={setChangeEntry} currentId={currentId}/>
        </div>
      </div>
    )
  }

  return(
  <>
  <div><CustomButton Text="Add new entry" onClickHandler={()=>{setAddNewEntry(true)}}/></div>
  {entries.map((entry: EntryProps) => <EntryBar 
    ingredients={entry.ingredients} 
    dish={entry.dish} 
    calories={entry.calories}
    fat={entry.fat}
    _id={entry._id}
    key={entry._id}
    reload = {reload}
    setReload={setReload}
    setChangeEntry={setChangeEntry}
    setCurrentId={setCurrentId}
    />
  )}
  </>
)

}

export default App;
