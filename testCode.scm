
(define x 1)
(define y 2)
(define z 32)

(define add
  (lambda (x y)
    (+ x y)))

(define fact
  (lambda (n)
    (if (> n 1)
        (* n (fact (- n 1)))
        1)))

(define (fact-tail n)
  (define fact-iter
    (lambda (n x)
      (if (> n 1)
          (fact-iter (- n 1) (* x n))
          x)))
  (fact-iter n 1))

(define fact-iter-let
  (lambda (n x)
    (let ((nn n)
          (xx x))
      (begin
        (if (> n 1)
            (fact-iter-let (- n 1) (* x n))
            x)))))

(add x y)
(fact-tail z)
