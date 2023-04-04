# arbitrary precision clebsch-gordan coefficients for SO(3)

## Content
* `gen-cg-table`: command line tool to generate CG table of two angular-momenta addition;

Example
```
./gen-cg-table ▶ go run main.go --j1=17/2 --j2=6
/var/folders/_0/2d8v_l8x5r947l5f35hdx0yw0000gq/T/clebsch-gordon.html
```

Then the CG table is rendered to an HTML indicated by the output line. See below for an example screenshot.


<img width="1534" alt="Screen Shot 2022-06-20 at 21 16 37" src="https://user-images.githubusercontent.com/107862003/174610115-af7bd8dd-5bbd-4e4f-9353-bceb1921de78.png">

* `./multi-angular`: command line tool to expand a tensor product of multiple angular momenta into the total angular momentum basis of the composite system.

Example
```
./multi-angular ▶ go run *.go --states "1/2,-1/2;3/2,1/2;1/2,1/2"
constructing C-G table for j1=3/2, j2=1/2 ...
constructing C-G table for j1=1, j2=1/2 ...
constructing C-G table for j1=2, j2=1/2 ...
/var/folders/_0/2d8v_l8x5r947l5f35hdx0yw0000gq/T/multi-angular.html
```

This takes input of 3 particles each in their respective spin state $|j_1,m_1\rangle=\left|\frac{1}{2},-\frac{1}{2}\right\rangle,|j_2,m_2\rangle=\left|\frac{3}{2},\frac{1}{2}\right\rangle,|j_3,m_3\rangle=\left|\frac{1}{2},\frac{1}{2}\right\rangle$, and calculates the product state's expansion into the total angular momentum basis $|j,m\rangle$ of the composite system.

The result is rendered to an HTML indicated by the output line. The last line of the page shows the desired expansion.

Note we have two distinct $\left|\frac{3}{2},\frac{1}{2}\right\rangle$ contributions from two disjoint irreducible 4-dimension subspaces $4_1$ and $4_2$ (indicated by the subscript).

<img width="1141" alt="Screenshot 2023-03-31 at 08 59 36" src="https://user-images.githubusercontent.com/107862003/228996535-857a5162-3c0a-4251-9341-d4771016adfe.png">

