import {ReactNode, useCallback, useEffect, useState} from "react";
import {
    Alert,
    AlertTitle, Avatar,
    Button,
    CircularProgress, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle,
    List, ListItemAvatar,
    ListItemButton,
    ListItemText,
    Snackbar,
    TextField,
    Typography
} from "@mui/material";
import HurtResultForm from "../hurtResult/HurtResultForm.tsx";
import {hurtNames, IAllResult, IItemInstance, IItemToSearch} from "../../interfaces.ts";
import InputComp from "./inputField/InputField.tsx";
import stratchInputType from "./inputField/handleInputTypes/stracher.ts";
import Box from '@mui/material/Box';
import {getHurtResult, getMultipleHurtResult} from "./resultGrabbers.ts";
import fetchWithAuth from "../../typeScriptFunc/fetchWithAuth.ts";
import {useNavigate} from "react-router-dom";


export default function MainSite() {
    // region zmienne
    const [Ean, setEan] = useState<string>("")
    const [prodName, setProdName] = useState<string>("")

    const [componentHashTable, setComponentHashTable] = useState<Map<hurtNames, ReactNode>>(new Map<hurtNames, ReactNode>())

    const [errorMessage, setErrorMessage] = useState<string>("")

    const [isLoadingProduct, setIsLoadingProduct] = useState<boolean>(false)


    const [prodToSearch, setProdToSearch] = useState<IItemToSearch[]>([])

    const [openSnackbar, setOpenSnackbar] = useState<boolean>(false)
    const [messageFromBackend, setMessageFromBackend] = useState<string>("")

    const [optItems, setOptItems] = useState<IItemInstance[]>([])
    const [allResult, setAllResult] = useState<IAllResult[]>([])

    const [open, setOpen] = useState<boolean>(false);
    const [agreement, setAgreement] = useState<boolean>(false)

    const navigate = useNavigate()

    const [fileName, setFileName] = useState<string>("")

    const [selectedIndex, setSelectedIndex] = useState(0);



    // endregion

    // region pozwala na przeciąganie plików
    const onDrop = useCallback((event: DragEvent) => {
        event.preventDefault();
        const file = event.dataTransfer?.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e: ProgressEvent<FileReader>) => {
                const result = e.target?.result;
                if (result && typeof result === 'string') {
                    setFileName(file.name);
                    setProdToSearch(stratchInputType(result, file.name)(result));
                }
            };

            reader.readAsText(file);
        }
    }, []);

    const onDragOver = useCallback((event: DragEvent) => {
        event.preventDefault();
    }, []);

    useEffect(() => {

        window.addEventListener('dragover', onDragOver);
        window.addEventListener('drop', onDrop);

        return () => {
            window.removeEventListener('dragover', onDragOver);
            window.removeEventListener('drop', onDrop);
        };
    }, [onDrop, onDragOver]);
    // endregion

    const changeResultComp = (ean: string, name: string) => {
        const newComponentHashTable = new Map<hurtNames, ReactNode>()
        setEan(ean)
        if (allResult.filter((item) => item.ean === ean).length === 0) {
            newComponentHashTable.set(hurtNames.none,
                <HurtResultForm
                    isCheaper={false}
                    name={hurtNames[hurtNames.none]}
                    priceForPack={-1}
                    princeForOne={-1}
                    productsInPack={-1}/>
            )
        } else {

            const eanRes = allResult.filter((item) => item.ean === ean)


            eanRes.map((newItem) => {
                const noEmpty = newItem.result.filter(i=>{
                    return i.Item.priceForOne!==-1
                })
                const cheaper = noEmpty.length!==0? noEmpty.reduce((p,n)=> {
                    return p.Item.priceForOne<n.Item.priceForOne? p : n
                }) : undefined
                newItem.result.map((newItem) => {

                    newComponentHashTable.set(newItem.hurtName,
                        <HurtResultForm
                            isCheaper={cheaper?.Item.hurtName===newItem.hurtName}
                            name={hurtNames[newItem.hurtName]}
                            priceForPack={newItem.Item.priceForPack}
                            princeForOne={newItem.Item.priceForOne}
                            productsInPack={newItem.Item.productsInPack}/>
                    )
                })
            })
        }
        setProdName(name)
        setComponentHashTable(newComponentHashTable)
    }

    const searchOneProd = () => {
        setIsLoadingProduct(true)
        setErrorMessage("")
        try {
            getHurtResult(Ean).then(data => {
                if (typeof (data) === "string") {
                    setProdName("brak Produktu")

                    if (data==="where logowanie?" || data=== "where Token?"){
                        navigate("/login")
                    }
                    setErrorMessage(data)

                    setIsLoadingProduct(false)
                    return
                }

                const newMap = new Map<hurtNames, ReactNode>()
                let i = 0

                const noEmpty = data.filter(i=>{
                    return i.priceForOne!==-1
                })
                const cheaper = noEmpty.length!==0? noEmpty.reduce((p,n)=> {
                    return p.priceForOne<n.priceForOne? p : n
                }) : undefined
                data.map((item) => {
                        setProdName(item.name)
                        i += 1
                        newMap.set(item.hurtName, (
                            <HurtResultForm
                                isCheaper={cheaper!==undefined? cheaper.hurtName===item.hurtName: false}
                                name={hurtNames[item.hurtName]}
                                priceForPack={item.priceForPack}
                                princeForOne={item.priceForOne}
                                productsInPack={item.productsInPack}
                            />
                        ))
                })
                if (i === 0) {
                    setProdName("brak produktu")
                    newMap.set(hurtNames.none, (
                        <HurtResultForm
                            isCheaper={false}
                            name={hurtNames[hurtNames.none]}
                            priceForPack={-1}
                            princeForOne={-1}
                            productsInPack={-1}
                        />
                    ))
                } else {
                    setComponentHashTable(newMap)
                }

                setIsLoadingProduct(false)
            });
        } catch (e: any) {
            setErrorMessage(e.message)
            setIsLoadingProduct(false)
        }
    }

    useEffect(() => {
        setErrorMessage("")
        if (prodToSearch.length === 0) {
            return;
        }

        setIsLoadingProduct(true)
        try {
            getMultipleHurtResult(prodToSearch).then(data => {

                if (typeof (data) == "string") { // jeżeli odpowiedzią jest string znaczy że mamy błąd
                    setErrorMessage(data)
                    setIsLoadingProduct(false)
                    return
                }

                const newOptItems: IItemInstance[] = []
                const newAllResult: IAllResult[] = []
                prodToSearch.map((item) => {
                    const ItemsMatchEan = data.get(item.Ean)
                    if (ItemsMatchEan) {
                        const optimalItem = ItemsMatchEan.reduce((prev, current) => {
                            // Check if the current price is -1; if so, skip it by returning prev
                            if (current.Item.priceForOne === -1) {
                                return prev;
                            }
                            // If prev price is -1 or current price is lower than prev price, select current
                            if (prev.Item.priceForOne === -1 || current.Item.priceForOne < prev.Item.priceForOne) {
                                return current;
                            }
                            // Otherwise, keep prev
                            return prev;
                        });
                        if (optimalItem.Item.priceForOne!==-1){
                            newOptItems.push({
                                name: item.Name,
                                ean: item.Ean,
                                item: optimalItem.Item,
                                count: item.Amount,
                            })
                        }
                        newAllResult.push({
                            ean: item.Ean,
                            result: ItemsMatchEan
                        })
                    }
                })
                setOptItems(newOptItems)
                setAllResult(newAllResult)
                setIsLoadingProduct(false)
                changeResultComp(prodToSearch[0].Ean,prodToSearch[0].Name)
            })
        } catch (e: any) {
            setErrorMessage(e.message)
            setIsLoadingProduct(false)
        }
        // zapisanie ich w optItems
    }, [prodToSearch])

    useEffect(() => {
        if (agreement && optItems.length!==0){
            setErrorMessage("")
            fetchWithAuth("/api/makeOrder", {
                method: "POST",
                body: JSON.stringify({Items: optItems.map(item => {
                        return {
                            Ean: item.ean,
                            Amount: item.count,
                            HurtName: item.item.hurtName
                        }
                    })}),
                headers: {
                    "Content-Type": "application/json",
                }

            }).then(response => {
                if (response.status !== 200) {
                    throw new Error("nie udało się złożyć zamówienia")
                }
                setOpenSnackbar(true)
                setMessageFromBackend("udało się dodać produkty do koszyka")
            }).catch(err => {
                setErrorMessage(err)
                throw new Error(err);
            })

        }
        setAgreement(false)

    }, [agreement]); //dodawanie do koszyka

    useEffect(() => {
        fetchWithAuth("/api/messages").then(response => {
            return response.json()
        }).then(data => {
            if (data.message == "") {
                return
            }
            setMessageFromBackend(data.message)
            setOpenSnackbar(true)
        })
    }, []) //pobieranie wiadomości


    return (
        <>
            <Box>
                <Box sx={{
                    padding: "20px",
                    display: "flex",
                    justifyContent: "center"
                }}>
                <Box className="photo" component="img" src={"./assets/logo.png"} alt="Logo" style={{
                    height: "150px"
                }}/>
                </Box>
            <p></p>
            <Box sx={{display: "flex", justifyContent: "space-around"}}>
                <Box sx={{width: "45%", display: "flex", justifyContent: "center"}}>
                    <TextField autoComplete={"off"} id="filled" label="skanuj pojedynczo"
                               placeholder="kod Ean" value={Ean}
                               onChange={e => setEan(e.target.value)} onKeyDown={e => {
                        if (e.key === "Enter") {
                            searchOneProd()
                        }
                    }}
                    />

                    <Button onClick={searchOneProd}>szukaj</Button>

                </Box>


                <TextField sx={{width: "45%"}} disabled autoComplete={"off"} id="filled-disabled"
                           label="nazwa produktu" value={prodName}/>

            </Box>
            <p></p>
            <Box style={{display: "flex", alignItems: "flex-start",justifyContent: "space-around"}}>
                {!isLoadingProduct ? (
                    <div className={"hurtResults"}
                         style={{width: "45%", display: "grid", gap: "10px"}}>
                        {
                            Array.from(componentHashTable.values()).map((element) => {
                                return element
                            })

                        }
                    </div>
                ) : <Box sx={{display: 'flex', padding: "20px"}}>
                    <CircularProgress/>
                </Box>
                }
                {
                    prodToSearch.length === 0 || isLoadingProduct ?
                        <></>
                        :
                        <List component="nav"
                              onKeyDown={(event) => {
                                  event.preventDefault()
                                  let newIndex=0
                                  if (event.key === 'ArrowUp') {
                                      newIndex = selectedIndex===0 ? 0 : selectedIndex - 1
                                  } else if (event.key === 'ArrowDown') {
                                      newIndex= selectedIndex===prodToSearch.length-1  ? prodToSearch.length-1 : selectedIndex + 1
                                  }
                                  setSelectedIndex(newIndex);
                                  changeResultComp(prodToSearch[newIndex].Ean,prodToSearch[newIndex].Name)
                              }}
                              sx={{width: "45%",
                                  overflow: "scroll",
                                  maxHeight: "600px",
                                  overflowX: "hidden",
                                  border: "1px solid white",
                                  borderRadius: "10px",
                              }}
                        >
                            {prodToSearch.map((item,index ) => {
                                return (
                                    <ListItemButton
                                        selected={selectedIndex === index}
                                        onClick={() => {
                                            changeResultComp(item.Ean, item.Name)
                                            setSelectedIndex(index)
                                        }}>
                                        <ListItemAvatar>
                                            <Avatar>{index + 1}</Avatar>
                                        </ListItemAvatar>
                                        <ListItemText primary={item.Name}/>
                                    </ListItemButton>
                                )
                            })}
                        </List>
                }


            </Box>
            <Button sx={{margin: "20px", padding: "5px"}}
                    variant="outlined" color="error" onClick={() => {
                setProdToSearch([])
            }}>Wyczyść Listę</Button>
            {optItems.length !== 0 ?
                <Button variant="contained" color="success" sx={{margin: "20px", padding: "5px"}} onClick={() => {
                    setOpen(true)
                }}>
                    dodaj produkty do koszyków w hurtowniach
                </Button> : <></>}
            {errorMessage !== "" ?
                <Alert severity="error">
                    <AlertTitle>Error</AlertTitle>
                    {errorMessage}
                </Alert>
                : null}

            {fileName && <Box marginTop="1rem" width="100%">
                <Typography variant="h6" component="h2" gutterBottom>
                    {"przetwarzany plik := " + fileName}
                </Typography>
            </Box>}

            <InputComp setItem={prod => setProdToSearch(prod)} setName={name => setFileName(name)}/>
                <Dialog
                    open={open}
                    onClose={()=>{setOpen(false)}}
                    aria-labelledby="alert-dialog-title"
                    aria-describedby="alert-dialog-description"
                >
                    <DialogTitle id="alert-dialog-title">
                        {"Czy dodać produkty do koszyka?"}
                    </DialogTitle>
                    <DialogContent>
                        <DialogContentText id="alert-dialog-description">
                            czy wyrażasz zgodę na dodanie produktów do koszyka, spowoduje to usunięcie aktualnego koszyka
                        </DialogContentText>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={()=> {
                            setAgreement(false)
                            setOpen(false)}}
                        >nie zgadzam się</Button>
                        <Button onClick={()=>{
                            setAgreement(true)
                            setOpen(false)
                            }} autoFocus>
                            Zgoda
                        </Button>
                    </DialogActions>
                </Dialog>



            <Snackbar
                anchorOrigin={{vertical: 'top', horizontal: 'center'}}
                open={openSnackbar}
                onClose={() => {
                    setOpenSnackbar(false)
                    setMessageFromBackend("")
                }}
                message={messageFromBackend}
            />
            </Box>
        </>
    )
}



