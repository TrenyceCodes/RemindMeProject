import { useEffect, useState } from 'react'
import './App.css'
import axios, { AxiosError } from 'axios'

function App() {
  const [message, setMessage] = useState("")
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  function onSubmit() {
    axios.post("http://localhost:8080/user/register", {
      username: username,
      email: email,
      password: password
    }).then((response) => {
      console.log(response.data.message);
      console.log(response.data.data);
      let message: string = response.data.message;
      setMessage(message)
    }).catch((error: AxiosError) => {
      alert(error)
    })
  }
  

  return (
    <>
      <h1>Create User</h1>
      
      <form action="" onSubmit={onSubmit}>
        <input type="text" name="username" id="username" placeholder="Tom" value={username} onChange={(e) => setUsername(e.target.value)}/><br/>
        <input type="text" name="username" id="username" placeholder="Tom@gmail.com" value={email} onChange={(e) => setEmail(e.target.value)}/><br/>
        <input type="text" name="username" id="username" placeholder="123456743" value={password} onChange={(e) => setPassword(e.target.value)}/><br/>
        <input type="submit" value="Create Account" />
      </form>
      {message}
    </>
  )
}

export default App
