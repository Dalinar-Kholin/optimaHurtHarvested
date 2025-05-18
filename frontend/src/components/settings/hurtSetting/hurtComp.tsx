import {hurtNames} from "../../../interfaces.ts";
import {Alert, AlertTitle, Button, TextField} from "@mui/material";
import {useState} from "react";
import fetchWithAuth from "../../../typeScriptFunc/fetchWithAuth.ts";
import Box from "@mui/material/Box";
interface IHurtComp{
    name : hurtNames
    fn : (username: string, pass : string, name : hurtNames) => void
}


export default function HurtComp({name, fn} : IHurtComp){
    const [username, setUsername] = useState<string>("")
    const [password, setPassword] = useState<string>("")
    const [error, setError] = useState<string>("")

    const sendRequest =()=>{
        const body = {
            username : username,
            password: password,
            hurtName: name
        }
        fetchWithAuth("/api/checkCredentials",  {
            body: JSON.stringify(body),
            method: "POST",
            headers:{
                "Content-Type": "application/json"
            }
        }).then(response =>{
            if (response.status!=200){
                return;
            }

            // jeżeli dostaliśmy 200 oznacza że dane są prawidłowe i możemy jest ustawić
            return response.json()
        }).then((data => {
            if (data.error!= undefined){
                setError(data.error)
                return
            }
            setError("")
            fn(username, password, name)

        }))


            .catch(_err =>{
            setError("network error")

        })
    }

    const handleClick = (e : any)=> {
        if (e.key=="Enter"){
            sendRequest()
        }
    }


    return(
        <>
            <Box sx={{display: "flex", justifyContent: "center"}}>
                <TextField value={username}  autoComplete={"off"} label={"username"}
                           onKeyDown={(e)=>{handleClick(e)}}
                           onChange={(e)=>{
                    setUsername(e.target.value)
                }}></TextField>
                <TextField value={password} type="password"  autoComplete={"off"} label={"password"} onChange={(e)=>{
                    setPassword(e.target.value)
                }}
                        onKeyDown={(e) => handleClick(e)}
                ></TextField>
                <Button onClick={() => {
                    sendRequest()
                }}>dodaj do zapisania</Button>
                {error === "" ? <div></div> : <Alert severity="error">
                    <AlertTitle>Error</AlertTitle>
                    {error}
                </Alert>}
            </Box>
        </>
    )
}