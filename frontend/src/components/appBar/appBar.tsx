import {AppBar, Button, Container, Toolbar} from "@mui/material";
import Box from "@mui/material/Box";
import {useNavigate} from "react-router-dom";
import "./appBar.css"
import {PATH} from "../../interfaces.ts";



/*interface IAppBarCustomed{
    iconLink : string
}*/

export default function AppBarCustomed(/*{iconLink}:IAppBarCustomed*/){

    const navigate = useNavigate()

    const pages : PATH[] = ["strona główna"/*,"płatności",  "cennik"*/,"ustawienia", "kontakt"] // tutaj jakby co dodać płatności i cennik
    return(
        <>
            <AppBar position="sticky" id={"appBarComp"}>
                <Container maxWidth="xl">
                    <Toolbar disableGutters>
                        {/*<Box className="photo" component="img" src={"./assets/optimaLogo.png"} alt="Logo" style={{
            borderRadius: '50%',
            height:'50px',
            width:'50px',
    }}
     onClick={()=>{
                            navigate('/main');
                        }} />*/} {/*ICOM: miejsce na logo*/}
                        <Box sx={{  display: "flex", width: "100%", justifyContent : "space-between"}}>
                            <Box >
                                {pages.map((page) => (
                                <Button
                                    color="inherit"
                                    key={page}
                                    onClick={()=>{
                                        navigate('/' + page);
                                    }}
                                >
                                    {page}
                                </Button>
                            ))}
                            </Box>
                            <Box>
                                <Button
                                    color="inherit"
                                    key={"wyloguj"}
                                    onClick={()=>{
                                        localStorage.clear()
                                        navigate('/login');
                                    }}
                                >
                                    {"wyloguj"}
                                </Button>
                            </Box>

                        </Box>
                    </Toolbar>
                </Container>
            </AppBar>
        </>
    )
}