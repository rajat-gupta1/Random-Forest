# Project \#3: Your Choice!

**Due: Tuesday, May 23th at 11:59pm**

## Assignment

The final project gives you the opportunity to show me what you learned
in this course and to build your own parallel system. In particular, you
should think about implementing a parallel system in the domain you are
most comfortable in (data science, machine learning, computer graphics,
etc.). The system should solve a problem that can benefit from some form
of parallelization. Similar to how the performance of an image
processing system benefits from parallel data decomposition of an image.
If you are having trouble coming up with a problem for your system to
solve then consider the following:

-   [Embarrassingly Parallel
    Topics](https://en.wikipedia.org/wiki/Embarrassingly_parallel)
-   [Parallel
    Algorithms](https://en.wikipedia.org/wiki/Parallel_computing#Algorithmic_methods)

You are free to implement any parallel algorithm you like. However, you
are required to at least have the following features in your parallel
system:

-   An input/output component that allows the program to read in data or
    receive data in some way. The system will perform some
    computation(s) on this input and produce an output result.

-   A sequential implementation of the system. Make sure to provide a
    usage statement.

-   **Parallel System Implementation**: A work-stealing and
    work-balancing algorithm using a **unbounded-dequeue** implemented
    as linked-list (i.e., a chain of nodes similar to project \#1). You
    **must** use the definitions defined in `proj3/concurrent`.
    Specifically, you must implement these functions that return an
    `ExecutorService`

    ``` go
    // NewWorkStealingExecutor returns an ExecutorService that is implemented using the work-stealing algorithm. 
    //@param capacity - The number of goroutines in the pool 
    //@param threshold - The number of items that a goroutine in the pool can
    // grab from the executor in one time period. For example, if threshold = 10 
    // this means that a goroutine can grab 10 items from the executor all at
    // once to place into their local queue before grabbing more items. It's
    // not required that you use this parameter in your implementation. 
    func NewWorkStealingExecutor(capacity, threshold int) ExecutorService {
       .... 
    }

    // NewWorkBalancingExecutor returns an ExecutorService that is implemented using the work-balancing algorithm. 
    //@param capacity - The number of goroutines in the pool 
    //@param threshold - The number of items that a goroutine in the pool can 
    //grab from the executor in one time period. For example, if threshold = 10
    // this means that a goroutine can grab 10 items from the executor all at
    // once to place into their local queue before grabbing more items. It's
    // not required that you use this parameter in your implementation. 
     //@param thresholdBalance - The threshold used to know when to perform
     // balancing. Remember, if two local queues are to be balanced the
     // difference in the sizes of the queue must be greater than or equal to
     // thresholdBalance. You must use this parameter in your implementation. 
    func NewWorkBalancingExecutor(capacity, thresholdQueue, thresholdBalance int) ExecutorService {
       .... 
    }

    // NewUnBoundedDEQueue returns an empty UnBoundedDEQueue
    func NewUnBoundedDEQueue() DEQueue {
        ...
    }
    ```

    **You are not allowed to modify/add/delete anything from the
    interfaces/types provided in the \`\`proj3/concurrent\`\` package.**
    You **cannot** modify the type signatures in the above functions.
    Golang does not have the keyword *volatile*; therefore, you will
    need to either use `sync.Mutex` or `sync.atomics` or any other
    `sync` package object. I would recommend using `sync.Mutex`. You
    will need to adapt the array based version shown in class to a
    linked-list version. A few additional notes:

    1.  To handle the ABA problem, make sure to define an internal node
        structure that holds the actual item being held in the dequeue.
        Every time an item is inserted then a new node is created for
        that item. The ABA problem only occurs when reusing memory
        addresses so creating unique ones will stop this from happening.
        Thus, you will not need a stamp mechanism in this
        implementation. However, it does come at the cost of cache
        performance.
    2.  Inside the `pi.go` file, we provide some
        snippets of code that use an `ExecutorService`. This code is
        from a working implementation that uses `ExecutorService` from
        previous year's projects. The programs show how to use `Future`,
        `Callable` and `Runnable` within your application portion of the
        assignment.
    3.  You should be able to reuse code from both of these
        implementations. Make sure to think about this while
        implementing both solutions.
    4.  Place your working-stealing implementation in the
        `proj3/conurrent/stealing.go` file.
    5.  Place your working-balancing implementation in the
        `proj3/conurrent/balancing.go` file.
    6.  Place your unbounded-dequeue implementation in the
        `proj3/conurrent/unbounded.go` file.

-   Provide a detailed write-up and analysis of your system. For this
    assignment, this write-up is required to have more detail to explain
    your parallel implementations since we are not giving you a problem
    to solve. See the **System Write-up** section for more details.

-   Provide all the dataset files you used in your analysis portion of
    your write up. If these files are to big then you need to provide us
    a link so we can easily download them from an external source.

-   These points also include design points. You should think about the
    modularity of the system you are creating. Think about splitting
    your code into appropriate packages, when necessary.

-   **You must provide a script or specific commands that shows/produces
    the results of your system**. We need to be able to enter in a
    single command in the terminal window and it will run and produce
    the results of your system. Failing to provide a straight-forward
    way of executing your system that produces its result will result in
    **significant deductions** to your score. We prefer running a simple
    command line script (e.g., shell-script or python3 script). However,
    providing a few example cases of possible execution runs will be
    sufficient enough.

-   We should also be able to run specific versions of the system. There
    should be a option (e.g. via command line argument) to run the
    sequential version, or the various parallel versions. Please make
    sure to document this in your report or via the printing of a usage
    statement.

-   You are free to use any additional standard/third-party libraries as
    you wish. However, all the parallel work is **required** to be
    implemented by you.

-   There is a directory called `proj3` with a single `go.mod` file
    inside your repositories. Place all your work for project 3 inside
    this directory.

### System Write-up

In prior assignments, we provided you with the input files or data to
run experiments against a your system and provide an analysis of those
experiments. For this project, you will do the same with the exception
that you will produce the data needed for your experiments. In all, you
should do the following for the writeup:

-   Run experiments with data you generate for both the sequential and
    parallel versions. As with the data provided by prior assignments,
    the data should vary the granularity of your parallel system. For
    the parallel version, make sure you are running your experiments
    with at least producing work for `N` threads, where
    `N = {2,4,6,8,12}`. You can go lower/larger than those numbers based
    on the machine you are running your system on. You are not required
    to run project 3 on the Peanut cluster. You can run it on your local
    machine and base your `N` threads on the number of logical cores you
    have on your local machine. If you choose to run your system on your
    local machine then please state that in your report and the your
    machine specifications as well.
-   Produce speedup graph(s) for those data sets. You should have one
    speedup graph per parallel implementation you define in your system.
    This means either one or two speedup graphs.

Please submit a report (pdf document, text file, etc.) summarizing your
results from the experiments and the conclusions you draw from them.
Your report should include your plot(s) as specified above and a
self-contained report. That is, somebody should be able to read the
report alone and understand what code you developed, what experiments
you ran and how the data supports the conclusions you draw. The report
**must** also include the following:

-   Describe in **detailed** of your system and the problem it is trying
    to solve.
-   A description of how you implemented your parallel solutions.
-   Describe the challenges you faced while implementing the system.
    What aspects of the system might make it difficult to parallelize?
    In other words, what to you hope to learn by doing this assignment?
-   Specifications of the testing machine you ran your experiments on
    (i.e. Core Architecture (Intel/AMD), Number of cores, operating
    system, memory amount, etc.)
-   What are the **hotspots** (i.e., places where you can parallelize
    the algorithm) and **bottlenecks** (i.e., places where there is
    sequential code that cannot be parallelized) in your sequential
    program? Were you able to parallelize the hotspots and/or remove the
    bottlenecks in the parallel version?
-   What limited your speedup? Is it a lack of parallelism?
    (dependencies) Communication or synchronization overhead? As you try
    and answer these questions, we strongly prefer that you provide data
    and measurements to support your conclusions.
-   Compare and contrast the two parallel implementations. Are there
    differences in their speedups?

## Don't know What to Implement?

If you are unsure what to implement then by default you can reimplement
the image processing assignment using the work-balancing and
work-stealing. For the advance feature, you need to be creative and
think about how you could adapt the image processing to use a MapReduce
paradigm. It could be the case you implement a slightly different
version of the image processing project just for MapReduce.

**You cannot reimplement project 1 or other assignments**.

## Design, Style and Cleaning up

Before you submit your final solution, you should, remove

-   any `Printf` statements that you added for debugging purposes and
-   all in-line comments of the form: "YOUR CODE HERE" and "TODO ..."
-   Think about your function decomposition. No code duplication. This
    homework assignment is relatively small so this shouldn't be a major
    problem but could be in certain problems.

Go does not have a strict style guide. However, use your best judgment
from prior programming experience about style. Did you use good variable
names? Do you have any lines that are too long, etc.

As you clean up, you should periodically save your file and run your
code through the tests to make sure that you have not broken it in the
process.

## Grading

For this project, we grade as follows:
 - 50% Completeness. You get full marks if your code implements the required features without deadlocks or race conditions.
 - 20% Performance. You get full marks if your code scales, i.e. it shows a reasonable speedup when more threads are added. 
 - 20% Writeup. You get full marks if you describe the problem that you are solving well, and explain your observed performance and speedup.
 - 10% Design and Style.

## Submission

Before submitting, make sure you've added, committed, and pushed all
your code to GitHub. You must submit your final work through Gradescope
(linked from our Canvas site) in the "Project \#3" assignment page via
two ways,

1.  **Uploading from Github directly (recommended way)**: You can link
    your Github account to your Gradescope account and upload the
    correct repository based on the homework assignment. When you submit
    your homework, a pop window will appear. Click on "Github" and then
    "Connect to Github" to connect your Github account to Gradescope.
    Once you connect (you will only need to do this once), then you can
    select the repsotiory you wish to upload and the branch (which
    should always be "main" or "master") for this course.
2.  **Uploading via a Zip file**: You can also upload a zip file of the
    homework directory. Please make sure you upload the entire directory
    and keep the initial structure the **same** as the starter code;
    otherwise, you run the risk of not passing the automated tests.

As a reminder, for this assignment, there will be **no autograder** on
Gradescope. We will run the program the CS Peanut cluster and manually
enter in the grading into Gradescope. However, you **must still submit
your final commit to Gradescope**.

