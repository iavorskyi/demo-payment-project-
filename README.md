# demo-payment-project-

1. It's used postgres as database in this demo.
2. DB is deployed in the cloud and ready to use all time.
3. For production, offcorse, all moneyflow should be in one transaction(with possibility to rollback changes). go-pg supports this approuch, decition was don't complicate the code.
4. There is should be much more tests, but for demo, I think, it's ok.
