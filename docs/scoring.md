# misc notes about scoring

2.1637 fieldWeight(doc.description.analyzed:lager in 639) product of:

1.0000 tf (termFreq(doc.descriptin.analyzed:lager) = 1)
4.3274 idf(docFreq=51, maxdocs = 1449)
0.5000 fieldNorm (field=doc.description.analyzed, doc=639)


tf(t in d)   =  	 sqrt(frequency)
idf(t)  =  	 1 + log ( numDocs / (docFreq+1) )
fieldNorm = 1/sqrt(numterms

'v' version

'f' field_id - field definition

'i' term_bytes 0xff field_id - num docs using this term in this field

't' term_bytes 0xff field_id doc_id - term frequence in field in doc

'n' field_id doc_id - num terms in field in doc (norm)

'b' doc_id - term 0xff field_id pairs



queryNorm:


x = sum of squares of idf(t) for each term in query

queryNorm = 1/sqrt(x)
