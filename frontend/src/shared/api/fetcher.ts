export interface ProblemDetails {
    type: string
    title: string
    status: number
    detail?: string
    instance?: string
}

export class APIError extends Error {
    status: number
    problem: ProblemDetails

    constructor(status: number, problem: ProblemDetails) {
        super(problem.title)
        this.name = 'APIError'
        this.status = status
        this.problem = problem
    }
}

async function fetcher<T>(url: string, options: RequestInit = {}): Promise<T> {
    const response = await fetch(url, {
        ...options,
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers || {}),
        },
    })

    if (!response.ok) {
        const contentType = response.headers.get('content-type')
        if (contentType && contentType.includes('application/problem+json')) {
            const problem = (await response.json()) as ProblemDetails
            throw new APIError(response.status, problem)
        }
        throw new APIError(response.status, {
            type: 'about:blank',
            title: response.statusText,
            status: response.status,
        })
    }

    if (response.status === 204) {
        return {} as T
    }

    return response.json() as Promise<T>
}

export default fetcher
