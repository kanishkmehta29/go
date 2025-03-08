interface ButtonProps{
  Text: string;
  onClickHandler: ()=>void;
}

export default function CustomButton(props:ButtonProps) {
  return (
    <>
      <button onClick={props.onClickHandler} className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
        {props.Text}
      </button>
    </>
  );
}
