import { useState } from "react";
import axios from "axios";

export default function AddEntryForm({setAddNewEntry}: {setAddNewEntry: (value: boolean) => void}) {
    const [dish,setDish] = useState('')
    const [ingredients,setIngredients] = useState('')
    const [calories,setCalories] = useState(0.0)
    const [fat,setFat] = useState(0.0)

  return (
    <>
    <div className="w-full max-w-xs">
      <form className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <div className="mb-4">
        <label
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="dish"
            >
              dish
            </label>
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="dish"
              type="text"
              placeholder="dish"
              onChange={(e) => setDish(e.target.value)}
            />
          </div>
          <div className="mb-4">
            <label
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="ingredients"
            >
              ingredients
            </label>
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="ingredients"
              type="text"
              placeholder="ingredients"
              onChange={(e) => setIngredients(e.target.value)}
            />
          </div>
          <div className="mb-4">
            <label
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="calories"
            >
              calories
            </label>
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="calories"
              type="text"
              placeholder="calories"
              onChange={(e) => setCalories(parseFloat(e.target.value) || 0)}
            />
          </div>
          <div className="mb-4">
            <label
              className="block text-gray-700 text-sm font-bold mb-2"
              htmlFor="fat"
            >
              fat
            </label>
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="fat"
              type="text"
              placeholder="fat"
              onChange={(e) => setFat(parseFloat(e.target.value) || 0)}
            />
          </div>
          <div className="flex items-center justify-between">
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
              type="button"
              onClick={async()=>{
                await axios.post("http://localhost:8080/entry/create",{
                    dish:dish,
                    ingredients:ingredients,
                    calories:calories,
                    fat:fat
                })
                setAddNewEntry(false)
              }}
            >
              Add
            </button>
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
              type="button"
              onClick={()=>{
                setAddNewEntry(false)
              }}
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </>
  );
}
