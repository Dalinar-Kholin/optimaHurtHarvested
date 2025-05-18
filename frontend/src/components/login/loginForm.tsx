import {useState} from "react";
import {Alert, AlertTitle, Button, TextField} from "@mui/material";
import "./login.css"
import {useNavigate} from "react-router-dom";
import {hurtNames} from "../../interfaces.ts";
import Box from "@mui/material/Box";



interface hurtLoginSuccess{
    hurt: number
    success : boolean
}

interface loginResponse{
    result : hurtLoginSuccess[]
    token : string
    availableHurts : number
    accountStatus : number
    companyName: string
}

interface badResponse{
    error : string
}

export default function LoginForm(){


    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [isProperData, setIsProperData] = useState<boolean>(true)
    const [errorMessage, setErrorMessage] = useState<string>("")

    const navigate = useNavigate()

/*
    const CustomTextField = styled(TextField)({
        '& input:-webkit-autofill': {
            '-webkit-box-shadow': '35 7a 38 1000px white inset !important',
        },
    });*/

    return (
        <>
            <Box sx={{
                padding: "50px",
                display: "flex",
                justifyContent: "center"
            }}>
                <Box className="photo" component="img" src={"./assets/logo.png"} alt="Logo" style={{
                    height: "300px"
                }}/>
            </Box>
            <form onSubmit={e => {
                e.preventDefault()
                // logowanie sie
                const LoginData = {
                    username: username,
                    password: password,
                }
                fetch('/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(LoginData),
                }).then(response => {
                    if (!response.ok) {
                        // Próbujemy pobrać dane błędu jako JSON
                        return response.json().then(errorData => {
                            setErrorMessage(errorData.error || "błąd sieci")
                            setIsProperData(false)
                            throw new Error('Error: ' + errorData.error);
                        });
                    }
                    return response.json();
                }).then((data: loginResponse | badResponse) => {

                    if ('error' in data) {
                        setErrorMessage(data.error);
                        setIsProperData(false);
                        return;
                    }

                    // Jeżeli `data` jest typu `loginResponse`, wykonujemy poniższe operacje
                    localStorage.setItem("accessToken", data.token);
                    localStorage.setItem("availableHurts", data.availableHurts.toString());
                    localStorage.setItem("accountStatus", data.accountStatus.toString());
                    localStorage.setItem("companyName", data.companyName);

                    // Przetwarzanie wyników logowania do hurtowni
                    data.result.map(i => {
                        if (!i.success) {
                            alert(`Nie udało się zalogować do hurtowni ${hurtNames[i.hurt]}`);
                        }
                    });

                    // Ustawiamy stan po pomyślnym logowaniu
                    setIsProperData(true);
                    setUsername("");
                    setPassword("");

                    // Przenoszenie użytkownika na stronę główną
                    navigate("/strona główna");

                }).catch(error => {
                    console.error('There has been a problem with your fetch operation:', error);
                })

            }}>
                <TextField
                    id="filled"
                    label="nazaw użytkownia"
                    placeholder="nazaw użytkownia"
                    value={username}
                    onChange={e => setUsername(e.target.value)}
                />
                <p></p>
                <TextField
                    id="outlined-password-input"
                    label="hasło"
                    type="password"
                    autoComplete="current-password"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                />
                <p></p>
                {isProperData ? <div></div> : <Alert severity="error">
                    <AlertTitle>Error</AlertTitle>
                    {errorMessage}
                </Alert>}
                <p></p>
                <Button variant="contained" type={"submit"}>
                    zaloguj
                </Button>
                <p></p>
                <Button variant="contained" onClick={() => {
                    navigate("/signIn")
                }}>
                    załóż konto
                </Button>
                <p></p>
                <Button variant="outlined" onClick={() => {
                    navigate("/forgotPassword")
                }}>zapomniałem hasła</Button>

            </form>
        </>
    )

}