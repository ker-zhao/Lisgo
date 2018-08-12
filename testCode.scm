
(define x 1)
(define y 2)
(define z 4)

(define add
  (lambda (x y)
    (+ x y)))

(define fact
  (lambda (n)
    (if (> n 1)
        (* n (fact (- n 1)))
        1)))

(add x y)
(fact z)
'