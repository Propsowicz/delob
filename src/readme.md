PLAN
1. OK - integration tests
2. OK - refactor tokenizer -> token scan -> parser
3. OK - add tcp
4. OK - add SCRAM (user, credentials, store, handshake, session store)
5. OK - add pipeline
6. OK - dockerized
7. backup
8. move tests to higher level
9. add atomicity to saving logs
FIX: order by elo asc; works but should throw since it needs to be Elo not elo
FIX: parse -> should not be able to use key more than once
