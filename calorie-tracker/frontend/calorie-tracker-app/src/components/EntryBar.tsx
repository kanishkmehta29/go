import CustomButton from "./CustomButton";
import axios from "axios";

export interface EntryProps{
  _id:string;
  dish: string;
  ingredients: string[];
  calories: number;
  fat: number;
  reload: boolean;
  setReload: React.Dispatch<React.SetStateAction<boolean>>;
  setCurrentId: React.Dispatch<React.SetStateAction<string>>;
  setChangeEntry: React.Dispatch<React.SetStateAction<boolean>>;
}

function EntryBar(props:EntryProps) {
  return (
    <>
      <div className="grid grid-cols-8 border border-gray-500 rounded-sm m-2">
        <div>Dish: {props.dish}</div>
        <div className="col-span-2">Ingredients: {props.ingredients}</div>
        <div>Calories: {props.calories}</div>
        <div>Fat: {props.fat}</div>
        <div><CustomButton Text="delete Entry" onClickHandler={async()=>{
          await axios.delete("http://localhost:8080/entry/delete/"+props._id)
          props.setReload(!props.reload)
        }}/></div>
        <div><CustomButton Text="change ingredients" onClickHandler={()=>{}}/></div>
        <div className="justify-self-end"><CustomButton Text="change entry" onClickHandler={()=>{
          props.setCurrentId(props._id)
          props.setChangeEntry(true)
        }}/></div>
      </div>
    </>
  );
}

export default EntryBar;
