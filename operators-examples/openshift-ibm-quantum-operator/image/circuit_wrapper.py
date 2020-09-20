class CircuitWrapper:

    circuitName = None
    qubits = None
    circuit = None
    status = None
    jobId = None
    backend = None
    answer = None

    def __init__(self, circuitName, qubits, circuit):
        self.circuitName = circuitName
        self.qubits = qubits
        self.circuit = circuit
