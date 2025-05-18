import {IHurtInfoForComp} from "./handleResult/handleResultInterfaces.ts";
import {handleResults} from "./handleResult/handleResults.ts";
import {hurtNames, IItemToSearch} from "../../interfaces.ts";
import fetchWithAuth from "../../typeScriptFunc/fetchWithAuth.ts";


export function getHurtResult(Ean: string): Promise<IHurtInfoForComp[] | string>  {
    const url = "/api/takePrice?" + new URLSearchParams({ean: Ean});


    let newData: IHurtInfoForComp[] = [];

    return fetchWithAuth(url, {
        credentials: "include",
        method: "GET",
    }).then(response => {
        if (!response.ok){
            throw response
        }

        return response.json();
    }).then(data => {

        if (data.error != undefined){
            return data.error
        }

        data.forEach((element: any) => {
            newData.push( handleResults({name: element.hurtName})(element.result));
        });


        return newData;
    }).catch(err => {
        if (err instanceof Response) {
            // Obsługa odpowiedzi z błędem
            return err.json().then(errorData => {
                return errorData.error
            }).catch(parseError => {
                return parseError
            });
        } else {
            // Inne typy błędów, np. brak połączenia z siecią
            return "błąd połączenia"
        }
    });
}

export interface IServerMultipleDataResult{
    Ean : string,
    Item : IHurtInfoForComp
    hurtName : hurtNames
}



export async function getMultipleHurtResult(Items: IItemToSearch[]):  Promise<Map<string, IServerMultipleDataResult[]> | string>{
    const map = new Map<string, IServerMultipleDataResult[]>();

    let isOk= true;

    const data = await fetchWithAuth("/api/takePrices", {
        credentials: "include",
        method: "POST",
        body: JSON.stringify({Items: Items}),
        headers: {
            "Content-Type": "application/json"
        }
    }).then(response => {
        if (!response.ok){
            throw response
        }
        return response.json();
    }).catch(err => {
        isOk= false
        if (err instanceof Response) {
            // Obsługa odpowiedzi z błędem
            return err.json().then(errorData => {
                return errorData.error
            }).catch(parseError => {
                return parseError
            });
        } else {
            // Inne typy błędów, np. brak połączenia z siecią
            return "błąd połączenia"
        }
    });
    if (!isOk){
        return data // zwracamy tutaj data, ponieważ w catch zwróciliśmy tam request
    }

    try {

        if (data.error !== undefined){
            return data.error
        }

        data.map((i : any) => { // mapujemy po konkretnych hurtowniach
            i.Result.map((item : any) => { // mapujemy po konkretnych wynikach z hurtowni
                const itemArray = map.get(item.Ean); // poprzednie wyniki
                const newItem = {
                    Ean: item.Ean,
                    Item: handleResults({name: i.HurtName})(item.Item), // przetwarzamy dane
                    hurtName: i.HurtName
                };

                if (itemArray !== undefined) {
                    itemArray.push(newItem); // jeżeli mamy poprzednie wyniki to pushujemy nasz nowy
                } else {
                    map.set(item.Ean, [newItem]); // jeżeli nie było wcześniej wyników to tworzymy nową tablicę z konkretnymi wynikamy
                }
            });
        });

        return map;
    } catch (err : any) {
        throw new Error(err.message);
    }
}
