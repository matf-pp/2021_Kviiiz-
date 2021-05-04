# 2021_Kviiiz-

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5e0b674ea3cb439e8079e0b2a901ebe7)](https://app.codacy.com/gh/matf-pp/2021_Kviiiz-?utm_source=github.com&utm_medium=referral&utm_content=matf-pp/2021_Kviiiz-&utm_campaign=Badge_Grade_Settings)

## Tema
Implementacija kviza zasnovana na chat room-u.

## Uputstvo
Korisnik postavlja svoje ime komandom **/name** *ime*.  
Nakon toga može da se priključi željenoj sobi komandom **/join** *room_name*. Ako je uneto ime sobe koja ne postoji napraviće se nova soba i korisnik postaje njen admin.  
Komandom **/rooms** se izlistavaju sve slobodne sobe.  
Članovi mogu da se dopisuju dok igra nije u toku.  
Admin pokreće kviz komandom **/start**, nakon čega se pojavljuje prvo pitanje. Kada svi takmičari daju svoje odgovore prelazi se na naredno pitanje, a igrači koji su tačno odgovorili dobijaju odgovarajući broj poena. Svako pitanje ima vremensko ograničenje nakon kog se prelazi na naredno pitanje, ukoliko svi odgovore pre isteka vremena odmah će se preći na naredno pitanje.
Jedan kviz se sastoji od 10 nasumično izabranih pitanja.  
Korisnik u bilo kom trenutku može uneti komandu **/help** da bi video koje opcije može da koristi.  
Soba se napušta komandom **/quit**.

## Implementacija
* jezik GO
* neke od korišćenih biblioteka: fmt, net, log, string, bufio, math, time 
* pitanja za kviz se preuzimaju sa stranice https://opentdb.com/

## Pokretanje

Pokretanje servera:
```
./Releases/2021_Kviiiz-
```
Pokretanje klijenta: 
```
telnet localhost 8888
```
## Operativni sistem
Linux
 
## Autori
* Milica Gnjatović (milicagnjatovic18@gmail.com)
* Slobodan Jenko (jenko.slobodan@gmail.com)
* Isidora Slavković (isidora.slavkovic@gmail.com)
