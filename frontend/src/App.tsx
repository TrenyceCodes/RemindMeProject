import { useEffect, useState } from 'react'
import './App.css'
import axios, { AxiosError } from 'axios'

function App() {
  const [message, setMessage] = useState("")

  useEffect(() => {
    axios.get("http://localhost:8080/").then((response) => {
      console.log(response.data.message);
      let message: string = response.data.message;
      setMessage(message)
    }).catch((error: AxiosError) => {
      alert(error)
    })
  }, [])
  

  return (
    <>
      <h1>hello</h1>
      {message}
    </>
  )
}

export default App
